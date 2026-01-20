package com.hunesion.webfluxv1.repository;

import com.hunesion.webfluxv1.model.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.r2dbc.repository.Modifying;
import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.LocalDateTime;
import java.util.UUID;

@Repository
public interface UserRepository extends R2dbcRepository<User, UUID> {
    Flux<User> findAllBy(Pageable pageable);

    @Modifying
    @Query("INSERT INTO users (id, username, email) " +
            "VALUES (:id, :username, :email)")
    Mono<User> insertUser(UUID id, String username, String email);

}
