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
