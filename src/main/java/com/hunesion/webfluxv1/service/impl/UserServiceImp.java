package com.hunesion.webfluxv1.service.impl;

import com.hunesion.webfluxv1.model.entity.User;
import com.hunesion.webfluxv1.model.request.UserRequest;
import com.hunesion.webfluxv1.repository.UserRepository;
import com.hunesion.webfluxv1.service.UserService;
import com.hunesion.webfluxv1.utils.ApiResponseWithPaginationUtils;
import lombok.AllArgsConstructor;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.support.PageableExecutionUtils;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.transaction.reactive.TransactionalOperator;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.Instant;
import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class UserServiceImp implements UserService {

    private final UserRepository userRepository;
    private final TransactionalOperator transactionalOperator;

    @Override
    public Mono<ApiResponseWithPaginationUtils<User>> findAllUsers(int page, int size, String sortBy, String sortDirection) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.fromString(sortDirection), sortBy));
        return userRepository.findAllBy(pageable)
                // 2. Collect them into a List
                .collectList()
                // 3. Combine with the total count
                .zipWith(userRepository.count())
                .map(tuple -> {
                    // Current page users
                    List<User> users = tuple.getT1();
                    long total = tuple.getT2();
                    // 4. Create a Page with both data and pagination info
                    Page<User> pageResult = PageableExecutionUtils.getPage(users, pageable, () -> total);
                    return ApiResponseWithPaginationUtils.itemsAndPaginationResponse(pageResult);
                });
    }

    @Override
    public Mono<User> findUserById(UUID id) {
        return userRepository.findById(id);
    }

    @Override
    public Mono<User> saveUser(UserRequest user) {
        System.out.println(">>> 1. saveUser method called");

        User newUser = User.builder()
                .id(UUID.randomUUID())
                .username(user.username())
                .email(user.email())
                .build();

        System.out.println(">>> 2. User object created: " + newUser);

        return userRepository.insertUser(newUser.getId(), newUser.getUsername(), newUser.getEmail());

    }

    @Override
    public Mono<User> updateUser(UserRequest user, UUID id) {
        return userRepository.findById(id)
                .flatMap(existingUser -> {
                    existingUser.setUsername(user.username());
                    existingUser.setEmail(user.email());
                    return userRepository.save(existingUser);
                })
                .switchIfEmpty(Mono.error(new RuntimeException("User not found")));
    }

    @Override
    public Mono<Void> deleteUser(UUID id) {
        return null;
    }
}
