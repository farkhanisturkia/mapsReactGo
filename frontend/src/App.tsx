import { Routes, Route, useNavigate } from 'react-router-dom';
import Driver from './pages/Driver';
import Admin from './pages/Admin';
import './App.css';

function App() {
  const navigate = useNavigate();

  return (
    <div>
      <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
        <button onClick={() => navigate('/Driver')}>Driver</button>
        <button onClick={() => navigate('/Admin')}>Admin</button>
      </div>
      <Routes>
        <Route path="/Driver" element={<Driver />} />
        <Route path="/Admin" element={<Admin />} />
      </Routes>
    </div>
  );
}

export default App;
