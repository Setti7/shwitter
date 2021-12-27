import { Avatar } from "@mui/material";
import { FC } from "react";
import User from "../models/user";

const UserAvatar: FC<{ user: User }> = ({ user }) => {
  return <Avatar alt={user.name} sx={{ width: 36, height: 36 }} />;
};

export default UserAvatar;
