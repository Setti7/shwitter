import { Box, Divider, Typography } from "@mui/material";
import { FC, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import ApiError from "../models/errors/ApiError";
import { FriendOrFollower, FriendOrFollowerType } from "../models/user";
import { getFollowers, getFriends } from "../services/user";
import FriendOrFollowerCard from "./FriendOrFollowerCard";

interface Props {
  type: FriendOrFollowerType;
}

const FriendOrFollowerList: FC<Props> = ({ type }) => {
  const [fofs, setFofs] = useState<FriendOrFollower[]>([]);
  const [error, setError] = useState<String | undefined>();

  const { userId } = useParams();

  useEffect(() => {
    const getData = async () => {
      if (userId) {
        let result: FriendOrFollower[] | ApiError;

        if (type === FriendOrFollowerType.Friend) {
          result = await getFriends(userId);
        } else {
          result = await getFollowers(userId);
        }

        if (result instanceof ApiError) {
          setError(result.getFormattedStatus());
        } else {
          setFofs(result);
        }
      }
    };

    getData();
  }, [userId, type]);

  return (
    <>
      {error !== undefined ? (
        <Typography textAlign="center">{error}</Typography>
      ) : (
        <></>
      )}

      {fofs.map((f) => {
        return (
          <Box mb={1}>
            <FriendOrFollowerCard friendOrFollower={f} />
            <Divider sx={{ marginTop: 2 }} />
          </Box>
        );
      })}
    </>
  );
};

export default FriendOrFollowerList;
