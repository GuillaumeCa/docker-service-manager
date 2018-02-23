import React from 'react';

import Typography from 'material-ui/Typography';

import Card, { CardHeader, CardMedia, CardContent, CardActions } from 'material-ui/Card';
import Button from 'material-ui/Button';

import LogViewer from './LogViewer';

import { BASE_URL } from './config';

class Service extends React.Component {
  state = {
    services: [],
    serviceID: '',
    updateServices: false,
    error: false,
  }

  componentDidMount() {
    this.getServices();
  }

  getServices() {
    fetch(BASE_URL + '/service')
      .then(res => res.json())
      .then(data => {
        this.setState({ services: data });
      })
      .catch(err => {
        this.setState({
          error: true,
        })
      })
  }

  updateServices() {
    this.setState({ updateServices: true });
    fetch(BASE_URL + '/service/update', {
      method: 'PUT',
    }).then(res => {
      this.setState({ updateServices: false });
    })
  }

  getLog(serviceID) {
    this.setState({ error: false });
    if (this.reader) this.reader.cancel();
    fetch(BASE_URL + '/service/logs?id=' + serviceID).then(res => {
      if (!res.body) return;
      this.reader = res.body.getReader();
      const decoder = new TextDecoder("utf-8");
      let logs = [];
      const read = () => {
        this.reader.read().then(data => {
          if (!data.done) {
            const line = decoder.decode(data.value);
            logs = logs.concat(line.split('\n'));
            this.setState({ logs });
            if (this.logArea && this.state.autoScroll) {
              this.logArea.scrollTop = this.logArea.scrollHeight;
            }
            read();
          }
        })
      }
      read();
    }).catch(err => {
      this.setState({ error: true, logs: [] });
    })
  }

  showLog = (serviceID) => e => {
    this.setState({ serviceID });
  }

  render() {

    const { serviceID, updateServices, error } = this.state;
    return (
      <div>
        {
          error &&
          <Typography variant="title" color="secondary">
            Connection error
            </Typography>
        }
        {
          !error &&
          <div>
            <Button
              style={{ marginBottom: 20 }}
              variant="raised"
              color="secondary"
              disabled={updateServices}
              onClick={() => {
                this.updateServices();
              }}>{updateServices ? 'Updating services' : 'Update all'}</Button>
            {
              this.state.services.length === 0 &&
              <Typography variant="title" color="secondary">
                No service available
                </Typography>
            }
            {
              this.state.services.map(srv => {
                const imageName = srv.Spec.TaskTemplate.ContainerSpec.Image.split('@')[0];
                const replicas = srv.Spec.Mode.Replicated.Replicas;
                return (
                  <Card key={srv.ID} style={{ marginBottom: 20 }}>
                    <CardHeader
                      title={srv.Spec.Name}
                      subheader={`${imageName} - ${replicas} instance${replicas > 1 ? 's' : ''}`} />
                    <CardContent>
                      {
                        serviceID === srv.ID &&
                        <LogViewer serviceID={srv.ID} />
                      }
                    </CardContent>
                    <CardActions>
                      <Button
                        style={{ margin: '0 5px' }}
                        size="small"
                        color="primary"
                        onClick={this.showLog(srv.ID)}>
                        Show Logs
                    </Button>
                    </CardActions>
                  </Card>
                )
              })
            }
          </div>
        }
      </div>
    )
  }
}

export default Service;