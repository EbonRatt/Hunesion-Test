package com.hunesion.webfluxv1.service;

import com.hunesion.webfluxv1.model.entity.User;
import com.hunesion.webfluxv1.model.request.UserRequest;
import com.hunesion.webfluxv1.utils.ApiResponseWithPaginationUtils;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.UUID;

public interface UserService {

    Mono<ApiResponseWithPaginationUtils<User>> findAllUsers(int page, int size, String sortBy, String sortDirection);
    Mono<User> findUserById(UUID id);
    Mono<User> saveUser(UserRequest user);
    Mono<User> updateUser(UserRequest user, UUID id);
    Mono<Void> deleteUser(UUID id);
}
