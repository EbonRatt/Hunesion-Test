package com.hunesion.testdemor2dbc.repository;

import com.hunesion.testdemor2dbc.model.entity.ItemTag;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.data.repository.reactive.ReactiveSortingRepository;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface ItemTagRepository extends R2dbcRepository<ItemTag,Long> {
    Flux<ItemTag> findAllByItemId(Long itemId);

    Mono<Integer> deleteAllByItemId(Long itemId);
}
