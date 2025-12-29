package com.hrd.testwebflux;

import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/users")
@RequiredArgsConstructor
@CrossOrigin(origins = "http://localhost:3000")
public class UserController {


    private final UserService service;

    @PostMapping
    public Mono<User> create(@RequestBody User user) {
        return service.create(user);
    }


    @GetMapping
    public Flux<User> findAll() {
        return service.findAll();
    }

//    @PutMapping("/{id}")
//    public Mono<User> update(@PathVariable Long id, @RequestBody User user) {
//        return service.update(id, user);
//    }


    @DeleteMapping("/{id}")
    public Mono<Void> delete(@PathVariable Long id) {
        return service.delete(id);
    }

}
