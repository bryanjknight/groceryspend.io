import './App.css';
import { Button } from './components/Button';
import axios from 'axios';
import { getBrowserInstance, EXTRACT_DOM_ACTION, ExtractDomResponse } from './browser';


const browser = getBrowserInstance();

export interface AppProps {
  webhookUrl: string;
}

export const App = (props: AppProps) => {

  const getAxios = (webhookUrl: string) => {
    return axios.create({
      baseURL: webhookUrl,
      // TODO: http proxy?
    })
  }

  const handleContentScriptResponse = (extractDomRequest: ExtractDomResponse) => {
    const payload = {
      url: window.location,
      timestamp: new Date().toISOString(),
      data: `<html><body>${extractDomRequest.dom}</body></html>`
    }
    const httpClient = getAxios(props.webhookUrl);

    httpClient.post('/', payload).then((resp) => {
      alert("Success");
    }).catch((err) => {
      alert(`Failure: ${err}`);
    })
  }

  const handleSendButtonClick = () => {

    // send a message to the content script to extract the dom
    browser.tabs.query({active: true}, (tabs) => {

      if (tabs.length !== 1) {
        console.error(`Expected 1 tab, got ${tabs.length}`);
        return;
      }
      const tab = tabs[0];
      if (!tab.id) {
        console.error(`Received null tab id`);
        return;
      }
      browser.tabs.sendMessage(tab.id, {action: EXTRACT_DOM_ACTION}, handleContentScriptResponse);
    })
    

  }

  return (
    <div className="App">
      <header className="App-header">
        <Button onClick={handleSendButtonClick} text={"Send to GrocerySpend.io"} />
      </header>
    </div>
  );
}
