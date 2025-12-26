package com.hunesion.testdemor2dbc.repository;

import com.hunesion.testdemor2dbc.model.entity.ItemTag;
import com.hunesion.testdemor2dbc.model.entity.Person;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface PersonRepository extends R2dbcRepository<Person,Long> {
}
