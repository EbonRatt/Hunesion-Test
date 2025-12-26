package com.hunesion.testdemor2dbc.exception;

public class UnexpectedItemVersionException extends RuntimeException {
    public UnexpectedItemVersionException(String message) {
        super(message);
    }

    public UnexpectedItemVersionException(Long expectedVersion, Long actualVersion) {
        super("Unexpected version: expected " + expectedVersion + ", but was " + actualVersion);
    }
}
