import React, { useState } from 'react';
import axios from 'axios';
import WeatherCard from './components/WeatherCard';
import './index.css';

const API_BASE_URL = (process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, '');

export default function App() {
  const [city, setCity] = useState('Hyderabad');
  const [weather, setWeather] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    const requestedCity = city.trim();
    if (!requestedCity) {
      setError('Enter a city name to continue.');
      setWeather(null);
      return;
    }
    setLoading(true);
    setError('');
    setWeather(null);
    try {
      const response = await axios.get(`${API_BASE_URL}/api/weather`, { params: { city: requestedCity }, timeout: 12000 });
      setWeather(response.data);
    } catch (requestError) {
      setError(requestError.response?.data?.error || 'Weather could not be retrieved. Check the city and try again.');
    } finally {
      setLoading(false);
    }
  }

  return <main className="page-shell">
    <section className="weather-panel" aria-labelledby="page-title">
      <div className="eyebrow">CITY WEATHER</div>
      <h1 id="page-title">Weather, at a glance.</h1>
      <p>Look up current conditions for a city around the world.</p>
      <form className="search-form" onSubmit={handleSubmit}>
        <label className="visually-hidden" htmlFor="city">City name</label>
        <input id="city" value={city} onChange={e => setCity(e.target.value)} placeholder="Enter a city, e.g. Hyderabad" />
        <button type="submit" disabled={loading}>{loading ? 'Checking...' : 'Get weather'}</button>
      </form>
      {error && <p className="message error" role="alert">{error}</p>}
      {loading && <p className="message" role="status">Contacting weather service...</p>}
      {weather && <WeatherCard weather={weather} />}
      <p className="footnote">Data provided by OpenWeatherMap. Temperatures are shown in Celsius.</p>
    </section>
  </main>;
}
