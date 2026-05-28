import axios from "axios";

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
