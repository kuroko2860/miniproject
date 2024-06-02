package com.kuroko.statistic.entity;

import org.locationtech.jts.geom.Point;
import jakarta.persistence.*;
import java.time.Instant;

@Entity
@Table(name = "objects")
public class ObjectEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "object_id", nullable = false)
    private String objectId;

    @Column(name = "type", nullable = false)
    private String type;

    @Column(name = "color", nullable = false)
    private String color;

    @Column(name = "location", columnDefinition = "GEOGRAPHY(POINT, 4326)", nullable = false)
    private Point location;

    @Column(name = "status", nullable = false)
    private String status;

    @Column(name = "timestamp", nullable = false)
    private Instant timestamp;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getObjectId() {
        return objectId;
    }

    public void setObjectId(String objectId) {
        this.objectId = objectId;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getColor() {
        return color;
    }

    public void setColor(String color) {
        this.color = color;
    }

    public Point getLocation() {
        return location;
    }

    public void setLocation(Point location) {
        this.location = location;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public Instant getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(Instant timestamp) {
        this.timestamp = timestamp;
    }

    // Getters and setters

}
