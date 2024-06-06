import { useState, useEffect } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css"; // Import Leaflet CSS
import "../index.css";

function Map() {
  const [objects, setObjects] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const position = [21.0278, 105.8342];

  useEffect(() => {
    const fetchObjects = async () => {
      try {
        const response = await fetch(
          "http://localhost:9090/api/objects/?start=2024-06-06T00:00:00&end=2024-06-07T00:00:00&longitude=105.8342&latitude=21.0278&distance=175"
        );
        if (!response.ok) {
          throw new Error("Network response was not ok.");
        }
        const data = await response.json();
        setObjects(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setIsLoading(false);
      }
    };

    fetchObjects();
  }, []);

  return (
    <div className="map">
      {isLoading ? (
        <p>Loading...</p>
      ) : error ? (
        <p>Error: {error}</p>
      ) : (
        <MapContainer center={position} zoom={13}>
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          {objects.map((object) => (
            <Marker
              key={object.id + object.createdAt}
              position={[object.location.latitude, object.location.longitude]} // Leaflet expects [lat, lng]
            >
              <Popup>
                <b>ID:</b> {object.id}
                <br />
                <b>Type:</b> {object.type}
                <br />
                <b>Color:</b> {object.color}
                <br />
                <b>Status:</b> {object.status}
              </Popup>
            </Marker>
          ))}
        </MapContainer>
      )}
    </div>
  );
}

export default Map;
