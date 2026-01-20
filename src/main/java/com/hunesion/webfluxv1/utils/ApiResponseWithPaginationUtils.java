package com.hunesion.webfluxv1.utils;

import com.hunesion.webfluxv1.model.response.PaginationResponse;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.domain.Page;

import java.util.List;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ApiResponseWithPaginationUtils<T> {
    private List<T> items;
    private PaginationResponse paginationResponse;


    public static <T> ApiResponseWithPaginationUtils<T> itemsAndPaginationResponse(Page<T> page) {
        return ApiResponseWithPaginationUtils.<T>builder()
                .items(page.getContent())
                .paginationResponse(PaginationResponse.paginationToResponse(page))
                .build();
    }
}
