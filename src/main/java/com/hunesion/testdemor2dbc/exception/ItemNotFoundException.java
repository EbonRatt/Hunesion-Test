package com.hunesion.testdemor2dbc.exception;

public class ItemNotFoundException extends RuntimeException {
    public ItemNotFoundException(String message) {
        super(message);
    }

    public ItemNotFoundException(Long id) {
        super("Item not found with id " + id);
    }
}
