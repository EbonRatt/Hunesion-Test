package com.hunesion.webfluxv1.model.request;

import com.hunesion.webfluxv1.model.enums.Rarity;

public record ItemRequest (
        String name,
        String description,
        double value,
        Rarity rarity
){}
