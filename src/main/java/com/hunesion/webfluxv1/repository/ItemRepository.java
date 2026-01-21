package com.hunesion.webfluxv1.repository;

import com.hunesion.webfluxv1.model.entity.Item;
import com.hunesion.webfluxv1.model.enums.Rarity;
import org.springframework.data.domain.Pageable;
import org.springframework.data.r2dbc.repository.Modifying;
import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.LocalDateTime;
import java.util.Collection;
import java.util.UUID;

@Repository
public interface ItemRepository extends R2dbcRepository<Item, UUID> {

    @Query("INSERT INTO items (id, name, description, value, rarity) " +
            "VALUES (:id, :name, :description, :value, :rarity) " +
            "RETURNING *")
    Mono<Item> insertItem(UUID id, String name, String description, double value, Rarity rarity);

    Flux<Item> findAllBy(Pageable pageable);

    @Modifying
    @Query("""
            DELETE FROM items
            WHERE id IN (:ids)
              AND (SELECT COUNT(*) FROM items WHERE id IN (:ids)) = :n
            """)
    Mono<Integer> deleteAllIfAllExist(Collection<UUID> ids, int n);

}
