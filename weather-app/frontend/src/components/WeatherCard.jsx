import React from 'react';

export default function WeatherCard({ weather }) {
  return <section className="result-card" aria-live="polite" aria-label={`Weather in ${weather.city}`}>
    <div className="result-heading">
      <div><p className="result-label">CURRENT CONDITIONS</p><h2>{weather.city}</h2><p className="condition">{weather.condition}</p></div>
      <div className="temperature">{Number(weather.temperature_c).toFixed(1)}<span>°C</span></div>
    </div>
    <div className="metrics">
      <div><span>Humidity</span><strong>{weather.humidity_percent}%</strong></div>
      <div><span>Wind speed</span><strong>{Number(weather.wind_speed_mps).toFixed(1)} m/s</strong></div>
    </div>
  </section>;
}
