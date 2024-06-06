package com.kuroko.statistic.repository;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.kuroko.statistic.entity.ObjectEntity;

@Repository
public interface ObjectRepository extends JpaRepository<ObjectEntity, Long> {

    @Query(value = "SELECT * FROM objects WHERE created_at >= :start AND created_at < :end AND ST_DWithin(location, ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326), :distance)", nativeQuery = true)
    List<ObjectEntity> findObjectsWithinDistanceAndTimeRange(
            @Param("start") LocalDateTime start,
            @Param("end") LocalDateTime end,
            @Param("longitude") double longitude,
            @Param("latitude") double latitude,
            @Param("distance") double distance);
}
