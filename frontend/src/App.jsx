import "leaflet-draw/dist/leaflet.draw.css"; // Import the Leaflet Draw CSS
import "leaflet/dist/leaflet.css";
import { useState } from "react";
import { FeatureGroup, MapContainer, TileLayer } from "react-leaflet";
import { EditControl } from "react-leaflet-draw";
import styles from "./App.module.css";

function App() {
  const [objectCount, setObjectCount] = useState(null);
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [typeFilter, setTypeFilter] = useState("all");
  const [colorFilter, setColorFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [mapLayers, setMapLayers] = useState([]);

  const handleStartTimeChange = (event) => setStartTime(event.target.value);
  const handleEndTimeChange = (event) => setEndTime(event.target.value);
  const handleTypeChange = (event) => setTypeFilter(event.target.value);
  const handleColorChange = (event) => setColorFilter(event.target.value);
  const handleStatusChange = (event) => setStatusFilter(event.target.value);
  const _onCreate = (e) => {
    const { layerType, layer } = e;
    if (layerType === "polygon") {
      const { _leaflet_id } = layer;

      setMapLayers((layers) => [
        ...layers,
        { id: _leaflet_id, latlngs: layer.getLatLngs()[0] },
      ]);
    }
  };

  const _onEdited = (e) => {
    const {
      layers: { _layers },
    } = e;

    Object.values(_layers).map(({ _leaflet_id, editing }) => {
      setMapLayers((layers) =>
        layers.map((l) =>
          l.id === _leaflet_id
            ? { ...l, latlngs: { ...editing.latlngs[0] } }
            : l
        )
      );
    });
  };

  const _onDeleted = (e) => {
    const {
      layers: { _layers },
    } = e;

    Object.values(_layers).map(({ _leaflet_id }) => {
      setMapLayers((layers) => layers.filter((l) => l.id !== _leaflet_id));
    });
  };

  const fetchData = async () => {
    try {
      const _polygon = [...mapLayers[0].latlngs, mapLayers[0].latlngs[0]]; // Close the polygon
      const polygonWKT = `POLYGON((${_polygon
        .map((point) => [point.lng, point.lat].join(" "))
        .join(", ")}))`;

      const response = await fetch(
        `http://localhost:9090/api/objects/count?start=${startTime}&end=${endTime}&polygon=${polygonWKT}&type=${typeFilter}&color=${colorFilter}&status=${statusFilter}`
      );

      if (response.ok) {
        const { object_count } = await response.json();
        setObjectCount(object_count);
      } else {
        throw new Error("Network response was not ok.");
      }
    } catch (err) {
      console.error("Error fetching object count:", err);
    }
  };

  return (
    <div className={styles.appContainer}>
      <div className={styles.sidebar}>
        <h2>Object Query</h2>

        <div className={styles.inputGroup}>
          <label htmlFor="start-time">Start Time:</label>
          <input
            type="datetime-local"
            id="start-time"
            value={startTime}
            onChange={handleStartTimeChange}
          />
        </div>

        <div className={styles.inputGroup}>
          <label htmlFor="end-time">End Time:</label>
          <input
            type="datetime-local"
            id="end-time"
            value={endTime}
            onChange={handleEndTimeChange}
          />
        </div>

        <div className={styles.inputGroup}>
          <label htmlFor="type">Type:</label>
          <select
            id="type"
            value={typeFilter}
            onChange={handleTypeChange}
            defaultValue={"all"}
          >
            <option value="all" defaultValue={true}>
              All
            </option>
            <option value="car">Car</option>
            <option value="truck">Truck</option>
            <option value="bike">Bike</option>
            <option value="bus">Bus</option>
          </select>
        </div>

        <div className={styles.inputGroup}>
          <label htmlFor="color">Color:</label>
          <select
            id="color"
            value={colorFilter}
            onChange={handleColorChange}
            defaultValue={"all"}
          >
            <option value="all" defaultValue={true}>
              All
            </option>
            <option value="red">Red</option>
            <option value="green">Green</option>
            <option value="blue">Blue</option>
            <option value="yellow">Yellow</option>
          </select>
        </div>

        <div className={styles.inputGroup}>
          <label htmlFor="status">Status:</label>
          <select
            id="status"
            value={statusFilter}
            defaultValue={"all"}
            onChange={handleStatusChange}
          >
            <option value="all" defaultValue={true}>
              All
            </option>
            <option value="moving">Moving</option>
            <option value="static">Static</option>
          </select>
        </div>

        <button className={styles.button} onClick={fetchData}>
          Query
        </button>

        {objectCount !== null && (
          <div className={styles.result}>
            <p>Object Count: {objectCount}</p>
          </div>
        )}
      </div>

      <MapContainer
        center={[21.0278, 105.8342]}
        zoom={13}
        style={{ height: "100vh", width: "100%" }}
      >
        <TileLayer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        />
        <FeatureGroup>
          <EditControl
            position="topright"
            draw={{
              rectangle: false,
              circle: false,
              marker: false,
              circlemarker: false,
              polyline: false,
            }}
            edit={{ remove: true }} // Allow to remove the shape
            onEdited={_onEdited}
            onCreated={_onCreate}
            onDeleted={_onDeleted}
          />
        </FeatureGroup>
      </MapContainer>
    </div>
  );
}

export default App;
