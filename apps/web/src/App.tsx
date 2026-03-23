import { useState } from 'react';

function App() {
  const [apiStatus, setApiStatus] = useState<string>('checking');

  React.useEffect(() => {
    // Check API health
    fetch('http://localhost:3001/health')
      .then((res) => res.json())
      .then((data) => {
        setApiStatus(data.status || 'unknown');
      })
      .catch(() => {
        setApiStatus('unreachable');
      });
  }, []);

  return (
    <div className="container" style={{ padding: '40px 20px' }}>
      <header style={{ marginBottom: '40px' }}>
        <h1 style={{ color: '#1976d2', marginBottom: '8px' }}>ML_Elec</h1>
        <p style={{ color: '#757575' }}>Plateforme de Maintenance Prédictive</p>
      </header>

      <main>
        <section className="card" style={{ marginBottom: '24px' }}>
          <h2 style={{ marginBottom: '16px' }}>👋 Bienvenue</h2>
          <p style={{ color: '#424242', marginBottom: '16px' }}>
            ML_Elec est en cours de développement. Cette interface permettra de :
          </p>
          <ul style={{ paddingLeft: '20px', color: '#424242' }}>
            <li>Acquérir des données depuis des capteurs (MQTT, OPC-UA)</li>
            <li>Détecter des anomalies avec des modèles de ML</li>
            <li>Prédire la durée de vie restante (RUL) des équipements</li>
            <li>Créer des pipelines de traitement no-code</li>
            <li>Visualiser les données en temps réel</li>
          </ul>
        </section>

        <section className="card">
          <h2 style={{ marginBottom: '16px' }}>🔧 État des Services</h2>
          <div style={{ display: 'flex', gap: '16px', flexWrap: 'wrap' }}>
            <div
              className="card"
              style={{
                flex: '1',
                minWidth: '200px',
                borderLeft: '4px solid',
                borderLeftColor: apiStatus === 'healthy' ? '#4caf50' : '#ff9800',
              }}
            >
              <h3 style={{ marginBottom: '8px' }}>Edge Agent API</h3>
              <p style={{ fontSize: '14px', color: '#757575' }}>Port 3001</p>
              <span
                style={{
                  display: 'inline-block',
                  marginTop: '8px',
                  padding: '4px 8px',
                  borderRadius: '4px',
                  fontSize: '12px',
                  fontWeight: '500',
                  backgroundColor: apiStatus === 'healthy' ? '#e8f5e9' : '#fff3e0',
                  color: apiStatus === 'healthy' ? '#2e7d32' : '#ef6c00',
                }}
              >
                {apiStatus === 'healthy'
                  ? '● En ligne'
                  : apiStatus === 'checking'
                  ? '○ Vérification...'
                  : '● Indisponible'}
              </span>
            </div>

            <div
              className="card"
              style={{
                flex: '1',
                minWidth: '200px',
                borderLeft: '4px solid #4caf50',
              }}
            >
              <h3 style={{ marginBottom: '8px' }}>Web Dashboard</h3>
              <p style={{ fontSize: '14px', color: '#757575' }}>Port 3000</p>
              <span
                style={{
                  display: 'inline-block',
                  marginTop: '8px',
                  padding: '4px 8px',
                  borderRadius: '4px',
                  fontSize: '12px',
                  fontWeight: '500',
                  backgroundColor: '#e8f5e9',
                  color: '#2e7d32',
                }}
              >
                ● En ligne
              </span>
            </div>

            <div
              className="card"
              style={{
                flex: '1',
                minWidth: '200px',
                borderLeft: '4px solid #4caf50',
              }}
            >
              <h3 style={{ marginBottom: '8px' }}>MQTT Broker</h3>
              <p style={{ fontSize: '14px', color: '#757575' }}>Port 1883</p>
              <span
                style={{
                  display: 'inline-block',
                  marginTop: '8px',
                  padding: '4px 8px',
                  borderRadius: '4px',
                  fontSize: '12px',
                  fontWeight: '500',
                  backgroundColor: '#e8f5e9',
                  color: '#2e7d32',
                }}
              >
                ● En ligne
              </span>
            </div>
          </div>
        </section>

        <section className="card" style={{ marginTop: '24px' }}>
          <h2 style={{ marginBottom: '16px' }}>📚 Documentation</h2>
          <p style={{ color: '#424242', marginBottom: '16px' }}>
            Consultez la documentation pour plus d'informations :
          </p>
          <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
            <a
              href="http://localhost:3001/health"
              target="_blank"
              rel="noopener noreferrer"
              className="btn btn-primary"
            >
              API Health Check
            </a>
            <a
              href="https://github.com/ml-elec/ml-elec"
              target="_blank"
              rel="noopener noreferrer"
              className="btn btn-primary"
            >
              GitHub
            </a>
          </div>
        </section>
      </main>

      <footer
        style={{
          marginTop: '40px',
          paddingTop: '20px',
          borderTop: '1px solid #e0e0e0',
          color: '#757575',
          fontSize: '14px',
        }}
      >
        <p>© 2026 ML_Elec. Open Source - License MIT</p>
      </footer>
    </div>
  );
}

export default App;
