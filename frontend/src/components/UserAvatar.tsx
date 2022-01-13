import { Avatar } from "@mui/material";
import { FC } from "react";
import User from "../models/user";
import { Link as RouterLink } from "react-router-dom";

const UserAvatar: FC<{ user: User }> = ({ user }) => {
  return (
    <RouterLink to={"/user/" + user.id}>
      <Avatar alt={user.name} sx={{ width: 36, height: 36 }} />
    </RouterLink>
  );
};

export default UserAvatar;
