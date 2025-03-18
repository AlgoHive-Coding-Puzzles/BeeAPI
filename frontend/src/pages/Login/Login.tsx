import { useState, useRef } from "react";
import { useNavigate, Navigate } from "react-router-dom";
import { Card } from "primereact/card";
import { InputText } from "primereact/inputtext";
import { Button } from "primereact/button";
import { Password } from "primereact/password";
import { Toast } from "primereact/toast";
import { useAuth } from "../../contexts/AuthContext";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const toast = useRef<Toast>(null);
  const navigate = useNavigate();
  const { login, isAuthenticated } = useAuth();

  // Redirect if already authenticated
  if (isAuthenticated) {
    return <Navigate to="/" />;
  }

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !password) {
      toast.current?.show({
        severity: "error",
        summary: "Error",
        detail: "Please enter username and password",
      });
      return;
    }

    try {
      setLoading(true);
      await login(username, password);
      navigate("/");
    } catch (error) {
      toast.current?.show({
        severity: "error",
        summary: "Login Failed",
        detail: error instanceof Error ? error.message : "Invalid credentials",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex justify-content-center align-items-center h-screen">
      <Toast ref={toast} />
      <Card title="Login" className="w-full mx-auto">
        <form onSubmit={handleLogin} className="p-fluid">
          <div className="field mb-4">
            <label htmlFor="username" className="block mb-2">
              Username
            </label>
            <InputText
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
              className="w-full"
            />
          </div>

          <div className="field mb-4">
            <label htmlFor="password" className="block mb-2">
              Password
            </label>
            <Password
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              toggleMask
              feedback={false}
              className="w-full"
            />
          </div>

          <Button
            label="Login"
            icon="pi pi-sign-in"
            type="submit"
            loading={loading}
            className="mb-3"
          />
        </form>
      </Card>
    </div>
  );
}
