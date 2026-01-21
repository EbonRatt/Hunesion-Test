package com.hunesion.webfluxv1.service.impl;

import com.hunesion.webfluxv1.mapper.ItemMapper;
import com.hunesion.webfluxv1.model.entity.Item;
import com.hunesion.webfluxv1.model.request.ItemRequest;
import com.hunesion.webfluxv1.model.response.ItemResponse;
import com.hunesion.webfluxv1.repository.ItemRepository;
import com.hunesion.webfluxv1.service.ItemService;
import com.hunesion.webfluxv1.utils.ApiResponseWithPaginationUtils;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.support.PageableExecutionUtils;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class ItemServiceImp implements ItemService {

    private final ItemRepository itemRepository;
    private final ItemMapper itemMapper;

    @Override
    public Mono<ItemResponse> createItem(ItemRequest itemRequest) {
        return itemRepository.insertItem(
                                UUID.randomUUID(),
                                itemRequest.name(),
                                itemRequest.description(),
                                itemRequest.value(),
                                itemRequest.rarity()
                        )
                        .map(itemMapper::itemToItemResponse).doOnSuccess(
                                itemResponse -> log.info("Item created: {}", itemResponse));
    }

    @Override
    public Mono<ApiResponseWithPaginationUtils<ItemResponse>> getAllItems(int page, int size, String sortBy, String sortDirection) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.fromString(sortDirection), sortBy));
        return itemRepository.findAllBy(pageable)
                .collectList()
                .zipWith(itemRepository.count())
                .map( tuple -> {
                    List<Item> items = tuple.getT1();
                    List<ItemResponse> itemResponses = items.stream().map(itemMapper::itemToItemResponse).collect(Collectors.toList());
                    long total = tuple.getT2();
                    Page<ItemResponse> pageResult = PageableExecutionUtils.getPage(itemResponses, pageable, () -> total);
                    return ApiResponseWithPaginationUtils.itemsAndPaginationResponse(pageResult);
                });
    }

    @Override
    public Mono<Void> deleteItems(List<UUID> uuidList) {
        List<UUID> ids = uuidList.stream().distinct().toList();
        int n = ids.size();
        if (n == 0) return Mono.empty();

        return itemRepository.deleteAllIfAllExist(ids, n)
                .flatMap(deleted -> {
                    if (deleted != n) {
                        return Mono.error(new RuntimeException("Some items not found"));
                    }
                    return Mono.empty();
                });
    }


}
