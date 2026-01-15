package com.hunesion.gucamole;

import org.apache.guacamole.GuacamoleException;
import org.apache.guacamole.io.GuacamoleReader;
import org.apache.guacamole.io.GuacamoleWriter;
import org.apache.guacamole.net.GuacamoleSocket;
import org.apache.guacamole.net.GuacamoleTunnel;
import org.apache.guacamole.protocol.GuacamoleInstruction;
import org.slf4j.Logger;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.time.LocalTime;

public class LoggingGuacamoleTunnel implements GuacamoleTunnel {

    private final GuacamoleTunnel wrappedTunnel;                      // The real tunnel
    private final Logger logger;                                     // For printing logs
    private final StringBuilder commandBuffer = new StringBuilder(); // Stores typed characters

    // Guacamole sends keyboard events as keysym codes (numbers). These are the codes we care about.
    private static final int KEYSYM_ENTER = 0xFF0D;
    private static final int KEYSYM_BACKSPACE = 0xFF08;
    private static final int KEYSYM_CTRL_U = 0x0015;  // Ctrl+U to clear line

    // Pattern: "3.key,<len>.<keysym>,<len>.<pressed>" e.g. "3.key,2.97,1.1"
    private static final Pattern KEY_PATTERN = Pattern.compile("\\d+\\.key,(\\d+)\\.(\\d+),(\\d+)\\.(\\d+)");

    // Blacklist file path - edit this file to change blocked commands in real-time
    private static final Path BLACKLIST_FILE = Paths.get("blacklist.txt");
    private static Set<String> blacklistedCommands = new HashSet<>();
    private static long lastModified = 0;

    // Time-based access control (must match GuacamoleController)
    private static final LocalTime ALLOWED_START = LocalTime.of(13, 36);   // 1:00 PM
    private static final LocalTime ALLOWED_END = LocalTime.of(13, 50);    // 1:20 PM

    /**
     * Load blacklist from file. Reloads automatically if file was modified.
     */
    private static synchronized Set<String> getBlacklistedCommands(Logger logger) {
        try {
            long currentModified = Files.getLastModifiedTime(BLACKLIST_FILE).toMillis();
            if (currentModified != lastModified) {
                Set<String> newBlacklist = new HashSet<>();
                for (String line : Files.readAllLines(BLACKLIST_FILE)) {
                    line = line.trim();
                    if (!line.isEmpty() && !line.startsWith("#")) {
                        newBlacklist.add(line.toLowerCase());
                    }
                }
                blacklistedCommands = newBlacklist;
                lastModified = currentModified;
                logger.info("Blacklist reloaded: {} commands", newBlacklist.size());
            }
        } catch (IOException e) {
            logger.warn("Could not read blacklist file: {}", e.getMessage());
        }
        return blacklistedCommands;
    }

    public LoggingGuacamoleTunnel(GuacamoleTunnel tunnel, Logger logger) {
        this.wrappedTunnel = tunnel;
        this.logger = logger;
    }

    @Override
    public UUID getUUID() {
        return wrappedTunnel.getUUID();
    }

    @Override
    public GuacamoleSocket getSocket() {
        return wrappedTunnel.getSocket();
    }

    @Override
    public GuacamoleReader acquireReader() {
        return wrappedTunnel.acquireReader();
    }

    @Override
    public void releaseReader() {
        wrappedTunnel.releaseReader();
    }

    @Override
    public boolean hasQueuedReaderThreads() {
        return wrappedTunnel.hasQueuedReaderThreads();
    }

    @Override
    public GuacamoleWriter acquireWriter() {
        GuacamoleWriter writer = wrappedTunnel.acquireWriter();
        // This returns our custom writer that intercepts keyboard input.
        return new LoggingGuacamoleWriter(writer, logger); // Return OUR wrapper!
    }

    @Override
    public void releaseWriter() {
        wrappedTunnel.releaseWriter();
    }

    @Override
    public boolean hasQueuedWriterThreads() {
        return wrappedTunnel.hasQueuedWriterThreads();
    }

    @Override
    public void close() throws GuacamoleException {
        wrappedTunnel.close();
    }

    @Override
    public boolean isOpen() {
        return wrappedTunnel.isOpen();
    }

    private class LoggingGuacamoleWriter implements GuacamoleWriter {

        private final GuacamoleWriter wrappedWriter;
        private final Logger logger;

        public LoggingGuacamoleWriter(GuacamoleWriter writer, Logger logger) {
            this.wrappedWriter = writer;
            this.logger = logger;
        }

        // This inner class intercepts everything going to guacd.
        @Override
        public void write(char[] message, int offset, int length) throws GuacamoleException {
            String raw = new String(message, offset, length);
            checkAndBlockCommand(raw);
            wrappedWriter.write(message, offset, length);
        }

        // This inner class intercepts everything going to guacd.
        @Override
        public void write(char[] message) throws GuacamoleException {
            String raw = new String(message);
            checkAndBlockCommand(raw);
            wrappedWriter.write(message);
        }

        @Override
        public void writeInstruction(GuacamoleInstruction instruction) throws GuacamoleException {
            if ("key".equals(instruction.getOpcode()) && instruction.getArgs().size() >= 2) {
                try {
                    int keysym = Integer.parseInt(instruction.getArgs().get(0)); // Which key?
                    int pressed = Integer.parseInt(instruction.getArgs().get(1)); // Press or release?
                    if (pressed == 1) { // Key pressed (not released)
                        // If command is blocked, checkKeyAndBlock returns true
                        // In that case, we already sent replacement and should NOT forward original Enter
                        if (checkKeyAndBlock(keysym)) {
                            return; // Don't forward this key to guacd
                        }
                    }
                } catch (NumberFormatException ignored) {
                }
            }
            wrappedWriter.writeInstruction(instruction); // Forward to guacd
        }

        private void checkAndBlockCommand(String raw) throws GuacamoleException {
            if (raw == null) return;
            Matcher m = KEY_PATTERN.matcher(raw);
            while (m.find()) {
                try {
                    int keysym = Integer.parseInt(m.group(2));
                    int pressed = Integer.parseInt(m.group(4));
                    if (pressed == 1) {
                        checkKeyAndBlock(keysym);
                    }
                } catch (NumberFormatException ignored) {
                }
            }
        }

        /**
         * Check if key should be blocked. Returns true if the key was handled (blocked),
         * false if it should be forwarded to guacd normally.
         */
        private boolean checkKeyAndBlock(int keysym) throws GuacamoleException {
            if (keysym == KEYSYM_ENTER) {
                // User pressed Enter! Check what they typed
                String command = commandBuffer.toString().trim();
                if (!command.isEmpty()) {
                    logger.info("==================================================");
                    logger.info("USER COMMAND: {}", command);
                    logger.info("==================================================");

                    // Check if current time is outside allowed hours
                    if (!isWithinAllowedTime()) {
                        logger.warn("BLOCKED: Time {} is outside allowed hours ({} - {})",
                                LocalTime.now(), ALLOWED_START, ALLOWED_END);
                        commandBuffer.setLength(0);
                        sendTimeBlockedMessage();
                        return true;
                    }

                    // Check if command is blacklisted
                    String baseCommand = command.split("\\s+")[0]; // Get First word, e.g., "ls" from "ls -la"
                    if (isBlacklisted(command) || isBlacklisted(baseCommand)) {
                        logger.warn("BLOCKED COMMAND: {} - This command is not allowed!", command);
                        commandBuffer.setLength(0); // Clear buffer

                        // Instead of throwing, send a replacement command
                        sendBlockedMessage(command);
                        return true; // Signal that we handled this key (don't forward Enter)
                    }
                }
                commandBuffer.setLength(0); // Clear buffer for next command
            } else if (keysym == KEYSYM_BACKSPACE) {

                // User pressed backspace - remove last character
                if (commandBuffer.length() > 0) {
                    commandBuffer.deleteCharAt(commandBuffer.length() - 1);
                }
            } else if (keysym >= 0x20 && keysym <= 0x7E) {
                // Printable character (a-z, 0-9, etc.) - add to buffer
                commandBuffer.append((char) keysym);
            }
            return false; // Not blocked, forward normally
        }

        /**
         * Display blocked message. Note: will show double prompt (limitation of terminal).
         */
        private void sendBlockedMessage(String blockedCommand) throws GuacamoleException {
            // Send Ctrl+U to clear current line
            sendKey(KEYSYM_CTRL_U, 1);
            sendKey(KEYSYM_CTRL_U, 0);

            // Use escape codes to overwrite the echo command line
            // Type printf command with ANSI escape codes
            String safeCommand = blockedCommand.replace("\"", "").replace("'", "");
            String message = "printf '\\e[2A\\e[2K\\e[31mBLOCKED COMMAND: " + safeCommand + " - This command is not allowed!\\e[0m\\n\\e[2K'";
            for (char c : message.toCharArray()) {
                sendKey(c, 1); // Press
                sendKey(c, 0); // Release
            }

            // Send Enter to execute
            sendKey(KEYSYM_ENTER, 1);
            sendKey(KEYSYM_ENTER, 0);
        }

        /**
         * Send a single key event to guacd.
         */
        private void sendKey(int keysym, int pressed) throws GuacamoleException {
            GuacamoleInstruction keyInstruction = new GuacamoleInstruction(
                    "key",
                    String.valueOf(keysym),
                    String.valueOf(pressed)
            );
            wrappedWriter.writeInstruction(keyInstruction);
        }

        private boolean isWithinAllowedTime() {
            LocalTime now = LocalTime.now();
            return !now.isBefore(ALLOWED_START) && !now.isAfter(ALLOWED_END);
        }

        private void sendTimeBlockedMessage() throws GuacamoleException {
            sendKey(KEYSYM_CTRL_U, 1);
            sendKey(KEYSYM_CTRL_U, 0);

            String message = "printf '\\e[2A\\e[2K\\e[31mACCESS DENIED: Current time is outside allowed hours (" + ALLOWED_START + " - " + ALLOWED_END + ")\\e[0m\\n\\e[2K'";
            for (char c : message.toCharArray()) {
                sendKey(c, 1);
                sendKey(c, 0);
            }
            sendKey(KEYSYM_ENTER, 1);
            sendKey(KEYSYM_ENTER, 0);
        }

        private boolean isBlacklisted(String command) {
            String lower = command.toLowerCase().trim();
            Set<String> blacklist = getBlacklistedCommands(logger);

            // Exact match: "ls" matches "ls"
            if (blacklist.contains(lower)) {
                return true;
            }

            // Prefix match: "ls -la" starts with "ls "
            for (String blocked : blacklist) {
                if (lower.equals(blocked) || lower.startsWith(blocked + " ")) {
                    return true;
                }
            }
            return false;
        }
    }
}
