package com.hunesion.webfluxv1.model.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.domain.Page;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class PaginationResponse {
    private Long totalElements;
    private Integer currentPage;
    private Integer pageSize;
    private Integer totalPages;

    public static  <T> PaginationResponse paginationToResponse(Page<T> page) {
        return PaginationResponse.builder()
                .currentPage(page.getNumber() + 1)
                .totalPages(page.getTotalPages())
                .totalElements(page.getTotalElements())
                .pageSize(page.getSize())
                .build();
    }
}
