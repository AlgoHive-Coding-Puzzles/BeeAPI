import { useState } from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import AppSidebar from "./components/Sidebar/Sidebar";
import HomePage from "./pages/Home/Home";
import ThemePage from "./pages/Theme/Theme";
import Login from "./pages/Login/Login";
import { AuthProvider, useAuth } from "./contexts/AuthContext";
import useFetch from "./hooks/useFetch";
import { ServerName } from "./types/ServerName";
import Forge from "./pages/Forge/Forge";

function AppContent() {
  const { data } = useFetch<ServerName>("/name");
  const [selectedMenu, setSelectedMenu] = useState("Home");
  const { isAuthenticated } = useAuth();

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  const renderContent = () => {
    switch (true) {
      case "Home" === selectedMenu:
        return <HomePage setSelectedMenu={setSelectedMenu} />;
      case selectedMenu.startsWith("Theme"):
        return <ThemePage selectedMenu={selectedMenu} />;
      case "Puzzles" === selectedMenu:
        return <p className="text-center">Get puzzle input</p>;
      case "Team" === selectedMenu:
        return <p className="text-center">Manage users</p>;
      case "Forge" === selectedMenu:
        return <Forge />;
      case "Settings" === selectedMenu:
        return <p className="text-center">Access Settings</p>;
    }
  };

  return (
    <div className="app-container">
      <AppSidebar
        selectedMenu={selectedMenu}
        setSelectedMenu={setSelectedMenu}
      />

      <div className="content-container">
        <h1 className="text-3xl font-bold text-center">{`BeeAPI - Interface - ${data?.name}`}</h1>
        <div className="mt-8">{renderContent()}</div>
      </div>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/*" element={<AppContent />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
