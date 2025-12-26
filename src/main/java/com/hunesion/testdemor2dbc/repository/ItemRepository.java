package com.hunesion.testdemor2dbc.repository;

import com.hunesion.testdemor2dbc.model.entity.Item;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ItemRepository extends R2dbcRepository<Item, Long> {
}
