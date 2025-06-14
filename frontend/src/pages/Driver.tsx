import { useState, useEffect } from 'react';

interface Point {
  name: string;
  lat: number;
  lng: number;
}

interface RouteRequest {
  current: Point;
  points: Point[];
}

function Driver() {
  const [route, setRoute] = useState<Point[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [points, setPoints] = useState<Point[]>([]);

  useEffect(() => {
    const fetchPoints = async () => {
      try {
        const response = await fetch('/data.json');
        if (!response.ok) {
          throw new Error(`Failed to fetch data.json: ${response.status}`);
        }
        const data: Point[] = await response.json();
        setPoints(data);
      } catch (error) {
        setError(error instanceof Error ? error.message : 'Unknown error fetching data.json');
      }
    };

    fetchPoints();
  }, []);

  const getNearestRoute = async () => {
    if (points.length === 0) {
      setError('No points available. Please check data.json.');
      return;
    }

    setIsLoading(true);
    setError(null);

    const data: RouteRequest = {
      current: { name: 'Current', lat: -6.2, lng: 106.8 },
      points: points,
    };

    try {
      const response = await fetch('https://mapsapi.msroot.my.id/api/route', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      setRoute(result.route);
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Unknown error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="App">
      <h1>Nearest Route Finder</h1>
      <button onClick={getNearestRoute} disabled={isLoading || points.length === 0}>
        {isLoading ? 'Loading...' : 'Get Route'}
      </button>
      {error && <p style={{ color: 'red' }}>Error: {error}</p>}
      {points.length === 0 && !error && <p>Loading points from data.json...</p>}
      {route.length > 0 && (
        <div>
          <h2>Route Result:</h2>
          <ol>
            {route.map((point, index) => (
              <li key={index}>
                {point.name} ({point.lat.toFixed(4)}, {point.lng.toFixed(4)})
              </li>
            ))}
          </ol>
        </div>
      )}
    </div>
  );
}

export default Driver;
