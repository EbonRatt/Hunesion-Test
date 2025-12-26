package com.hunesion.testdemor2dbc.service;

import lombok.Value;

@Value
public class ItemDeleted implements Event {

    Long itemId;

}
