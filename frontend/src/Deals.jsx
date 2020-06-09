import React, { useEffect, useState } from "react";
import Axios from "axios";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import Divider from "@material-ui/core/Divider";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import List from "@material-ui/core/List";
import CircularProgress from "@material-ui/core/CircularProgress";

function ListItemLink(props) {
  return <ListItem button component="a" {...props} />;
}

export default function Deals() {
  const [isLoading, setLoading] = useState(true);
  const [state, setState] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      await Axios.get("/api/v1/deals")
        .then((response) => {
          const data = [];
          response.data.forEach((item) => {
            data.push(item);
          });
          data.sort((a, b) => (a.votes < b.votes ? 1 : -1));
          setState(data);
          setLoading(false);
        })
        .catch(() => {});
    };
    fetchData();
  }, []);

  const getData = state.map((item) => {
    return (
      <div key={item.id}>
        <ListItemLink
          href={item.link}
          target="_blank"
          rel="noopener noreferrer"
        >
          <ListItemText primary={item.title} secondary={`+${item.votes}`} />
        </ListItemLink>
        <Divider variant="middle" />
      </div>
    );
  });

  return (
    <>
      <Typography variant="h4" component="h4" gutterBottom>
        Deals from the last 48 hours
      </Typography>
      <Paper elevation={24}>
        <List dense>
          {isLoading ? (
            <ListItem>
              <CircularProgress />
            </ListItem>
          ) : (
            <>{getData}</>
          )}
        </List>
      </Paper>
    </>
  );
}
