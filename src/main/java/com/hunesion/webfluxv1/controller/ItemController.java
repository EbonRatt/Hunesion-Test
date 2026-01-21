package com.hunesion.webfluxv1.controller;

import com.hunesion.webfluxv1.model.request.ItemRequest;
import com.hunesion.webfluxv1.model.response.ApiResponse;
import com.hunesion.webfluxv1.model.response.ItemResponse;
import com.hunesion.webfluxv1.service.ItemService;
import com.hunesion.webfluxv1.utils.ApiResponseWithPaginationUtils;
import com.hunesion.webfluxv1.utils.ResponseUtils;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/items")
@RequiredArgsConstructor
public class ItemController {

    private final ItemService itemService;

    @GetMapping
    public Mono<ResponseEntity<ApiResponse<ApiResponseWithPaginationUtils<ItemResponse>>>> getAllItems(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size,
            @RequestParam(defaultValue = "id") String sortBy,
            @RequestParam(defaultValue = "asc") String sortDirection
    ) {
        return ResponseUtils.responseSingle("Items fetched successfully", HttpStatus.OK, itemService.getAllItems(page, size, sortBy, sortDirection));
    }

    @PostMapping
    public Mono<ResponseEntity<ApiResponse<ItemResponse>>> createItem(@RequestBody ItemRequest itemRequest) {
        return ResponseUtils.responseSingle("Item created successfully", HttpStatus.OK, itemService.createItem(itemRequest));
    }

    @DeleteMapping
    public Mono<ResponseEntity<ApiResponse<Void>>> deleteItems(@RequestBody List<UUID> uuidList) {
        return ResponseUtils.responseSingle("Items deleted successfully", HttpStatus.OK, itemService.deleteItems(uuidList));
    }

}
