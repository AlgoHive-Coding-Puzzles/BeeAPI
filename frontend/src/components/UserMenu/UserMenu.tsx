import { Button } from "primereact/button";
import { useAuth } from "../../contexts/AuthContext";

export default function UserMenu() {
  const { username, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
    window.location.href = "/login";
  };

  return (
    <div className="user-menu p-3 mt-auto">
      <div className="text-center mb-3">
        <i className="pi pi-user text-2xl"></i>
        <div className="mt-2 text-sm font-bold">{username}</div>
      </div>

      <div className="ml-1">
        <Button icon="pi pi-sign-out" text onClick={handleLogout} />
      </div>
      <p>Logout</p>
    </div>
  );
}
