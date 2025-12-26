package com.hunesion.testdemor2dbc.service;

import com.hunesion.testdemor2dbc.model.dto.response.ItemResourceResponse;
import lombok.Value;

@Value
public class ItemSaved implements Event {

    ItemResourceResponse item;

}
