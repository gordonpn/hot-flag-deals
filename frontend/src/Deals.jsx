import { makeStyles } from "@material-ui/core/styles";
import {
  CircularProgress,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
  Typography,
} from "@material-ui/core";
import React, { useEffect, useState } from "react";
import Axios from "axios";

const useStyles = makeStyles((theme) => ({
  root: {
    width: "100%",
    // maxWidth: 360,
    backgroundColor: theme.palette.background.paper,
  },
}));

function ListItemLink(props) {
  return <ListItem button component="a" {...props} />;
}

export default function Deals() {
  const classes = useStyles();
  const [isLoading, setLoading] = useState(true);
  const [state, setState] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const result = await Axios.get("/api/v1/deals");
      const data = [];
      result.data.forEach((item) => {
        data.push(item);
      });
      data.sort((a, b) => (a.votes < b.votes ? 1 : -1));
      setState(data);
      setLoading(false);
    };
    fetchData();
  }, []);

  return (
    <>
      <Typography variant="h4" component="h4" gutterBottom>
        Deals from the last 48 hours
      </Typography>
      <Paper elevation={24}>
        {isLoading ? (
          <List dense>
            <ListItemIcon>
              <CircularProgress />
            </ListItemIcon>
          </List>
        ) : (
          <>
            <List dense>
              {state.map((item) => {
                return (
                  <div key={item.id}>
                    <ListItemLink
                      href={item.link}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <ListItemText
                        primary={item.title}
                        secondary={`+${item.votes}`}
                      />
                    </ListItemLink>
                  </div>
                );
              })}
            </List>
          </>
        )}
      </Paper>
    </>
  );
}
