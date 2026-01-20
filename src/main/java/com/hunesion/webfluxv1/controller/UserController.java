package com.hunesion.webfluxv1.controller;

import com.hunesion.webfluxv1.model.entity.User;
import com.hunesion.webfluxv1.model.request.UserRequest;
import com.hunesion.webfluxv1.model.response.ApiResponse;
import com.hunesion.webfluxv1.repository.UserRepository;
import com.hunesion.webfluxv1.service.UserService;
import com.hunesion.webfluxv1.utils.ApiResponseWithPaginationUtils;
import com.hunesion.webfluxv1.utils.ResponseUtils;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping
    public Mono<ResponseEntity<ApiResponse<ApiResponseWithPaginationUtils<User>>>> getAllUsers(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size,
            @RequestParam(defaultValue = "id") String sortBy,
            @RequestParam(defaultValue = "asc") String sortDirection
    ) {
        return ResponseUtils.responseSingle("Users fetched successfully", HttpStatus.OK, userService.findAllUsers(page, size, sortBy, sortDirection));
    }

    @GetMapping("/{id}")
    public Mono<ResponseEntity<ApiResponse<User>>> getUserById(@PathVariable UUID id) {
        return ResponseUtils.responseSingle("User fetched successfully", HttpStatus.OK, userService.findUserById(id));
    }

    @PostMapping
    public Mono<ResponseEntity<ApiResponse<User>>> saveUser(@RequestBody UserRequest user) {
        return ResponseUtils.responseSingle("User saved successfully", HttpStatus.OK, userService.saveUser(user));
    }

    @PutMapping("/{id}")
    public Mono<ResponseEntity<ApiResponse<User>>> updateUser(@PathVariable UUID id, @RequestBody UserRequest user) {
        return ResponseUtils.responseSingle("User updated successfully", HttpStatus.OK, userService.updateUser(user, id));
    }

    @DeleteMapping("/{id}")
    public Mono<ResponseEntity<ApiResponse<Void>>> deleteUser(@PathVariable UUID id) {
        return ResponseUtils.responseSingle("User deleted successfully", HttpStatus.OK, userService.deleteUser(id));
    }

}
