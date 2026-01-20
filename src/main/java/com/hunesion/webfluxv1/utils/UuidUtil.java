package com.hunesion.webfluxv1.utils;

import java.nio.ByteBuffer;
import java.util.UUID;

public class UuidUtil {
    public static byte[] toBytes(UUID uuid) {
        if (uuid == null) {
            return null;
        }
        ByteBuffer bb = ByteBuffer.wrap(new byte[16]);
        bb.putLong(uuid.getMostSignificantBits());
        bb.putLong(uuid.getLeastSignificantBits());
        return bb.array();
    }

    public static UUID fromBytes(byte[] bytes) {
        if (bytes == null || bytes.length != 16) {
            return null;
        }
        ByteBuffer bb = ByteBuffer.wrap(bytes);
        return new UUID(bb.getLong(), bb.getLong());
    }
}
