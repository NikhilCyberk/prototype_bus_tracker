// import { useState } from 'react'
// import reactLogo from './assets/react.svg'
// import viteLogo from '/vite.svg'
// import './App.css'

// function App() {
//   const [count, setCount] = useState(0)

//   return (
//     <>
//       <div>
//         <a href="https://vite.dev" target="_blank">
//           <img src={viteLogo} className="logo" alt="Vite logo" />
//         </a>
//         <a href="https://react.dev" target="_blank">
//           <img src={reactLogo} className="logo react" alt="React logo" />
//         </a>
//       </div>
//       <h1>Vite + React</h1>
//       <div className="card">
//         <button onClick={() => setCount((count) => count + 1)}>
//           count is {count}
//         </button>
//         <p>
//           Edit <code>src/App.jsx</code> and save to test HMR
//         </p>
//       </div>
//       <p className="read-the-docs">
//         Click on the Vite and React logos to learn more
//       </p>
//     </>
//   )
// }

// export default App

import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import Header from "./components/Header";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Dashboard from "./pages/Dashboard";
import RouteView from "./pages/RouteView";
import StopView from "./pages/StopView";
import { useAuth } from "./context/AuthContext";
import NotFound from "./pages/NotFound";

function App() {
  const { user } = useAuth();

  // Protected route component
  const ProtectedRoute = ({ children }) => {
    if (!user) {
      return <Navigate to="/login" />;
    }
    return children;
  };

  return (
    <>
      <Header />
      <main className="container py-4">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route path="/routes/:id" element={<RouteView />} />
          <Route path="/stops/:id" element={<StopView />} />
          // Add this to the existing Routes component
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      <footer className="bg-dark text-white text-center py-4">
        <div className="container">
          <p className="mb-0">
            CityBus Tracker &copy; {new Date().getFullYear()}. All rights
            reserved.
          </p>
        </div>
      </footer>
    </>
  );
}

export default App;
