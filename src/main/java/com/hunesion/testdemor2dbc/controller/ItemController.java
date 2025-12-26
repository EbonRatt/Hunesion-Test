package com.hunesion.testdemor2dbc.controller;

import com.hunesion.testdemor2dbc.exception.ItemNotFoundException;
import com.hunesion.testdemor2dbc.mapper.ItemMapper;
import com.hunesion.testdemor2dbc.model.dto.request.ItemPatchResource;
import com.hunesion.testdemor2dbc.model.dto.request.ItemUpdateResourceRequest;
import com.hunesion.testdemor2dbc.model.dto.request.NewItemResourceRequest;
import com.hunesion.testdemor2dbc.model.dto.response.ItemResourceResponse;
import com.hunesion.testdemor2dbc.service.Event;
import com.hunesion.testdemor2dbc.service.ItemDeleted;
import com.hunesion.testdemor2dbc.service.ItemSaved;
import com.hunesion.testdemor2dbc.service.impl.ItemServiceImp;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.Duration;

import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.linkTo;
import static org.springframework.http.ResponseEntity.created;
import static org.springframework.http.ResponseEntity.noContent;

@RestController
@RequestMapping("/items")
@RequiredArgsConstructor
@Slf4j
public class ItemController {

    private final ItemServiceImp itemService;
    private final ItemMapper itemMapper;


    @PostMapping
    public Mono<ResponseEntity<Void>> create(@Valid @RequestBody final NewItemResourceRequest newItemResource) {

        return itemService.create(itemMapper.toModel(newItemResource))
                .map(item -> created(linkTo(ItemController.class).slash(item.getId()).toUri()).build());
    }

    @GetMapping(value = "/getAll",produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<ItemResourceResponse> findAll() {
        return itemService.findAll().map(itemMapper::toResource);
    }

    @GetMapping(value = "/{id}")
    public Mono<ItemResourceResponse> findById(@PathVariable final Long id) {

        return itemService.findById(id, null, true).map(itemMapper::toResource);
    }

    @PutMapping(value = "/{id}")
    public Mono<ResponseEntity<Void>> update(@PathVariable @NotNull final Long id,
                                             @RequestHeader(value = HttpHeaders.IF_MATCH) final Long version,
                                             @Valid @RequestBody final ItemUpdateResourceRequest itemUpdateResource) {

        // Find the item and update the instance
        return itemService.findById(id, version, false)
                .map(item -> itemMapper.update(itemUpdateResource, item))
                .flatMap(itemService::update)
                .map(item -> noContent().build());
    }

//    @PatchMapping(value = "/{id}")
//    public Mono<ResponseEntity<Void>> patch(@PathVariable @NotNull final Long id,
//                                            @RequestHeader(value = HttpHeaders.IF_MATCH) final Long version,
//                                            @Valid @RequestBody final ItemPatchResource patch) {
//
//        return itemService.findById(id, version, true)
//                .map(item -> itemMapper.patch(patch, item))
//                .flatMap(itemService::update)
//                .map(itemId -> noContent().build());
//    }
//
//    private Mono<Boolean> verifyExistence(Long id) {
//        return itemRepository.existsById(id).handle((exists, sink) -> {
//            if (Boolean.FALSE.equals(exists)) {
//                sink.error(new ItemNotFoundException(id));
//            } else {
//                sink.next(exists);
//            }
//        });
//    }

    @DeleteMapping("/{id}")
    public Mono<ResponseEntity<Void>> delete(@PathVariable final Long id,
                                             @RequestHeader(value = HttpHeaders.IF_MATCH) final Long version) {

        return itemService.deleteById(id, version)
                .map(empty -> noContent().build());
    }

    @GetMapping("/events")
    public Flux<ServerSentEvent<Event>> listenToEvents() {

        final Flux<Event> itemSavedFlux =
                this.itemService.listenToSavedItems()
                        .map(itemMapper::toResource)
                        .map(ItemSaved::new);

        final Flux<Event> itemDeletedFlux =
                this.itemService.listenToDeletedItems()
                        .map(ItemDeleted::new);

        return Flux.merge(itemSavedFlux, itemDeletedFlux)
                .map(event -> ServerSentEvent.<Event>builder()
                        .retry(Duration.ofSeconds(4L))
                        .event(event.getClass().getSimpleName())
                        .data(event).build());
    }
}
