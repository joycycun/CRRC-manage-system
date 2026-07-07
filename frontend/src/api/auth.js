import axios from "axios";
import request from "@/utils/request";

export const loginApi = (username, password) => {
  return axios.post(
    "/api/login",
    { username, password },
    {
      headers: {
        "Content-Type": "application/json"
      }
    }
  );
};

export const changePasswordApi = data => {
  return request.post("/change-password", data);
};

export const getRoleOptionsApi = () => {
  return request.get("/roles");
};

export const createUserApi = data => {
  return request.post("/users", data);
};

export const getUsersApi = () => {
  return request.get("/users");
};

export const deleteUserApi = id => {
  return request.delete(`/users/${id}`);
};

export const getUserPermissionsApi = id => {
  return request.get(`/users/${id}/permissions`);
};

export const saveUserPermissionsApi = (id, permissions) => {
  return request.post(`/users/${id}/permissions`, { permissions });
};
