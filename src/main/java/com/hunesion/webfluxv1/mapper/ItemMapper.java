package com.hunesion.webfluxv1.mapper;

import com.hunesion.webfluxv1.model.entity.Item;
import com.hunesion.webfluxv1.model.request.ItemRequest;
import com.hunesion.webfluxv1.model.response.ItemResponse;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface ItemMapper {
    Item toEntity(ItemRequest itemRequest);
    ItemResponse itemToItemResponse(Item item);
}
