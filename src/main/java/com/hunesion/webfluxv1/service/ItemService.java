package com.hunesion.webfluxv1.service;

import com.hunesion.webfluxv1.model.request.ItemRequest;
import com.hunesion.webfluxv1.model.response.ItemResponse;
import com.hunesion.webfluxv1.utils.ApiResponseWithPaginationUtils;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.UUID;

public interface ItemService {
    Flux<ItemResponse> createItems(List<ItemRequest> itemRequest);
    Mono<ApiResponseWithPaginationUtils<ItemResponse>> getAllItems(int page, int size, String sortBy, String sortDirection);
    Mono<Void> deleteItems(List<UUID> uuidList);
}
