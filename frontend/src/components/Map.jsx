// src/components/MapComponent.js

import { useEffect, useState } from "react";
import { MapContainer, TileLayer, Polygon, Marker, Popup } from "react-leaflet";
import axios from "axios";
import "leaflet/dist/leaflet.css";

const Map = () => {
  const [objects, setObjects] = useState([]);
  const [polygon, setPolygon] = useState([
    [51.505, -0.09],
    [51.51, -0.1],
    [51.51, -0.12],
  ]); // Example coordinates
  setPolygon([
    [51.505, -0.09],
    [51.51, -0.1],
    [51.51, -0.12],
  ]);
  useEffect(() => {
    const fetchData = async () => {
      const response = await axios.get("http://localhost:9090/api/objects", {
        params: {
          coordinates: polygon.map(([lat, lng]) => [lng, lat]), // Reverse to match Leaflet's [lng, lat] format
          startTime: "2024-06-01T00:00:00Z",
          endTime: "2024-06-15T00:00:00Z",
        },
      });
      setObjects(response.data);
    };

    fetchData();
  }, [polygon]);

  return (
    <MapContainer
      center={[51.505, -0.09]}
      zoom={13}
      style={{ height: "100vh", width: "100%" }}
    >
      <TileLayer
        url="https://tile.openstreetmap.org/13/555/123.png"
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      />
      <Polygon positions={polygon} />
      {objects.map((obj) => (
        <Marker
          key={obj.id}
          position={[obj.location.coordinates[1], obj.location.coordinates[0]]}
        >
          <Popup>
            <div>
              <p>Type: {obj.type}</p>
              <p>Color: {obj.color}</p>
              <p>Status: {obj.status}</p>
              <p>Timestamp: {new Date(obj.timestamp).toLocaleString()}</p>
            </div>
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
};

export default Map;
