package com.kuroko.statistic.repository;

import org.locationtech.jts.geom.Polygon;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.kuroko.statistic.entity.ObjectEntity;

import java.time.Instant;
import java.util.List;

@Repository
public interface ObjectRepository extends JpaRepository<ObjectEntity, Long> {

    @Query("SELECT o FROM ObjectEntity o WHERE o.location && :polygon AND o.timestamp BETWEEN :startTime AND :endTime AND (:type IS NULL OR o.type = :type) AND (:color IS NULL OR o.color = :color)")
    List<ObjectEntity> findByGeoAndTimeRangeAndAttributes(
            @Param("polygon") Polygon polygon,
            @Param("startTime") Instant startTime,
            @Param("endTime") Instant endTime,
            @Param("type") String type,
            @Param("color") String color);
}
