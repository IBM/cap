import java.io.ByteArrayInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.ArrayList;

import javax.xml.bind.JAXBContext;
import javax.xml.bind.JAXBElement;

import org.apache.http.client.HttpClient;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.BasicResponseHandler;
import org.apache.http.impl.client.DefaultHttpClient;

import com.sun.xml.internal.ws.util.Pool.Marshaller;
import com.sun.xml.internal.ws.util.Pool.Unmarshaller;

import oasis.names.tc.emergency.cap._1.Alert;

import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.DocumentBuilder;
import org.w3c.dom.Document;
import org.w3c.dom.NodeList;
import org.w3c.dom.Node;
import org.w3c.dom.Element;

public class CapFeed {
	final static String NwsNationalAtomFeedURL = "https://alerts.weather.gov/cap/us.php?x=1";

	@SuppressWarnings("deprecation")
	public static final HttpClient client = new DefaultHttpClient();

	public static void main(String[] args) {
		// TODO Auto-generated method stub
		Document doc = getFeed();
		System.out.println(getAlertsAll(doc));
	}

	public static Document getFeed() {
		// NwsNationalAtomFeedURL is the URL for the NWS National Atom feed
		HttpGet method = new HttpGet(NwsNationalAtomFeedURL);
		ResponseHandler<String> handler = new BasicResponseHandler();
		String responseStr;
		try {
			responseStr = client.execute(method, handler);

			FileOutputStream out = new FileOutputStream("alerts.txt");
			out.write(responseStr.getBytes());
			out.close();

			DocumentBuilderFactory dbFactory = DocumentBuilderFactory.newInstance();
			DocumentBuilder dBuilder = dbFactory.newDocumentBuilder();
			Document doc = dBuilder.parse(new ByteArrayInputStream(responseStr.getBytes()));
			doc.getDocumentElement().normalize();
			return doc;
		} catch (Exception e) {
			e.printStackTrace();
		}
		return null;
	}

	public static ArrayList<Alert> getAlertsAll(Document doc) {
		ArrayList<Alert> alerts = new ArrayList<Alert>();

		try {
			NodeList nList = doc.getElementsByTagName("entry");
			System.out.println("Number of entry's:" + nList.getLength());
			for (int j = 0; j < nList.getLength(); j++) {
				NodeList children = nList.item(j).getChildNodes();
				for (int i = 0; i < children.getLength(); i++) {
					Node child = children.item(i);
					if (child.getNodeName().equalsIgnoreCase("link")) {
						alerts.add(processLink(child.getAttributes().getNamedItem("href").getNodeValue()));
					}
				}
			}
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		System.out.println("Bye");
		return alerts;
	}

	public static Alert processLink(String url) throws Exception {
		HttpGet method = new HttpGet(url);
		ResponseHandler<String> handler = new BasicResponseHandler();
		String responseStr;

		JAXBContext jc = JAXBContext.newInstance("oasis.names.tc.emergency.cap._1");
		javax.xml.bind.Unmarshaller unmarshaller = jc.createUnmarshaller();

		responseStr = client.execute(method, handler);

		@SuppressWarnings("unchecked")
		Alert feed = (Alert) unmarshaller.unmarshal(new ByteArrayInputStream(responseStr.getBytes()));
		javax.xml.bind.Marshaller marshaller = jc.createMarshaller();
		marshaller.setProperty(javax.xml.bind.Marshaller.JAXB_FORMATTED_OUTPUT, true);
		marshaller.marshal(feed, System.out);

		return feed;
	}

	public static ArrayList<Alert> alertByType(ArrayList<Alert> list, String msgType) {
		ArrayList<Alert> alerts = new ArrayList<Alert>();

		for (Alert alert : list) {
			if (alert.getMsgType().equalsIgnoreCase(msgType)) {
				alerts.add(alert);
			}
		}
		return alerts;
	}
}
