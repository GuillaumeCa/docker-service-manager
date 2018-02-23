import React from 'react';

import {
  FormControlLabel,
} from 'material-ui/Form';
import Checkbox from 'material-ui/Checkbox';
import { BASE_URL_WS } from './config';

class LogViewer extends React.Component {


  state = {
    logs: [],
    autoScroll: true,
  }


  componentDidMount() {
    if (this.props.serviceID) {
      this.initWS(this.props.serviceID);
    }
  }

  componentWillUnmount() {
    if (this.ws) {
      this.ws.close();
    }
  }

  initWS(serviceID) {
    this.ws = new WebSocket(BASE_URL_WS + '/ws');
    this.ws.onopen = () => {
      this.ws.send(JSON.stringify({
        type: 'SERVICE_LOG',
        data: serviceID
      }));
    }
    this.ws.onmessage = (msgEvent) => {
      try {
        const msg = JSON.parse(msgEvent.data)
        this.handleLog(msg)
      } catch (error) {
        console.error("could not read message: ", msgEvent, error.message);
      }
    }
  }

  handleLog(msg) {
    if (msg.Content) {
      const logs = [...this.state.logs, msg.Content];
      this.setState({ logs });
      this.autoScroll();
    }
  }

  onAutoScroll = () => {
    this.setState({ autoScroll: !this.state.autoScroll });
  }

  autoScroll = () => {
    if (this.logArea && this.state.autoScroll) {
      this.logArea.scrollTop = this.logArea.scrollHeight;
    }
  }


  render() {
    return (
      <div>
        <div style={{ marginBottom: 20 }}>
          <FormControlLabel
            control={
              <Checkbox
                checked={this.state.autoScroll}
                onChange={this.onAutoScroll}
                value="autoscroll"
              />
            }
            label="Autoscroll"
          />
          {this.state.logs.length === 0 && <p>No log to display</p>}
        </div>
        {
          !this.state.error && this.state.logs.length > 0 &&
          <div ref={ref => this.logArea = ref} style={{
            fontFamily: 'monospace',
            height: 300,
            overflow: 'auto',
            borderRadius: 3,
            background: 'black',
            color: 'white',
            fontSize: 14,
            fontWeight: 'bold',
            lineHeight: 1.2,
            padding: 15,
          }}>
            {this.state.logs.map((l, i) => <p style={{ margin: 0 }} key={i}>{l}</p>)}
          </div>
        }
      </div>
    )
  }
}

export default LogViewer;