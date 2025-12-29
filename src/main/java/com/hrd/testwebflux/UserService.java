package com.hrd.testwebflux;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
@RequiredArgsConstructor
@Slf4j
public class UserService {

    private final UserRepository repository;
    private final UserEventPublisher publisher;


    @Transactional(transactionManager = "connectionFactoryTransactionManager")
    public Mono<User> create(User user) {
        return repository.save(user)
                .doOnNext(u -> log.info("Saved user with ID: {}", u.getId()))
                .flatMap(u -> repository.findById(u.getId())
                        .doOnNext(found -> log.info("Verified in DB: {}", found))
                        .switchIfEmpty(Mono.defer(() -> {
                            log.error("User not found after save!");
                            return Mono.empty();
                        })))
                .flatMap(u -> {
                    publisher.publish(new UserEvent("USER_CREATED", u));
                    return Mono.just(u);
                });
    }


    public Flux<User> findAll() {
        return repository.findAll();
    }


//    public Mono<User> update(Long id, User user) {
//        return repository.findById(id)
//                .flatMap(existing -> {
//                    existing.setName(user.getName());
//                    existing.setEmail(user.getEmail());
//                    return repository.save(existing);
//                })
//                .doOnSuccess(u -> publisher.publish("User updated: " + u.getId()));
//    }


    public Mono<Void> delete(Long id) {
        return repository.deleteById(id)
                .doOnSuccess(v -> publisher.publish(new UserEvent("USER_DELETED", id)));
    }

}
