import { Box, Divider, Typography } from "@mui/material";
import { FC, useEffect, useState } from "react";
import ApiError from "../models/errors/ApiError";
import { Timeline } from "../models/shweet";
import { getUserline } from "../services/shweets";
import ShweetCard from "./ShweetCard";

interface Props {
  userId: string;
}

const UserLine: FC<Props> = ({ userId }) => {
  const [userline, setUserline] = useState<Timeline>([]);
  const [error, setError] = useState<String | undefined>();

  useEffect(() => {
    const getData = async () => {
      const lineResult = await getUserline(userId);

      if (lineResult instanceof ApiError) {
        setError(lineResult.getFormattedStatus());
      } else {
        setUserline(lineResult);
      }
    };

    getData();
  }, [userId]);

  return (
    <>
      {error !== undefined ? (
        <Typography textAlign="center">{error}</Typography>
      ) : (
        <></>
      )}

      {userline.map((s) => {
        return (
          <Box mb={1}>
            <ShweetCard initialShweet={s} />
            <Divider sx={{ marginTop: 2 }} />
          </Box>
        );
      })}
    </>
  );
};

export default UserLine;
