/**
 * Tests E2E pour ML_Elec Web Dashboard
 * Scénario: Dashboard principal et vérification des services
 */

import { test, expect } from '@playwright/test';

test.describe('ML_Elec Dashboard', () => {
  test('should load the dashboard successfully', async ({ page }) => {
    // Naviguer vers le dashboard
    await page.goto('/');

    // Vérifier le titre
    await expect(page).toHaveTitle(/ML_Elec/);

    // Vérifier l'en-tête
    const heading = page.getByRole('heading', { name: /ML_Elec/ });
    await expect(heading).toBeVisible();
  });

  test('should display service status cards', async ({ page }) => {
    await page.goto('/');

    // Vérifier les cartes de statut
    const edgeAgentCard = page.getByText(/Edge Agent API/);
    await expect(edgeAgentCard).toBeVisible();

    const webDashboardCard = page.getByText(/Web Dashboard/);
    await expect(webDashboardCard).toBeVisible();

    const mqttBrokerCard = page.getByText(/MQTT Broker/);
    await expect(mqttBrokerCard).toBeVisible();
  });

  test('should check API health endpoint', async ({ page }) => {
    // Faire une requête directe à l'API health
    const response = await page.request.get('http://localhost:3001/health');

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.status).toBe('healthy');
    expect(data.timestamp).toBeDefined();
  });
});

test.describe('Authentication Flow', () => {
  test('should display login form', async ({ page }) => {
    await page.goto('/');

    // Chercher le lien ou bouton de login
    // Note: À implémenter quand l'UI de login sera créée
    test.skip();
  });

  test('should authenticate with valid credentials', async ({ page }) => {
    test.skip('Login UI not implemented yet');

    await page.goto('/login');

    // Remplir le formulaire
    await page.getByLabel('Username').fill('admin');
    await page.getByLabel('Password').fill('admin');
    await page.getByRole('button', { name: 'Login' }).click();

    // Vérifier la redirection
    await expect(page).toHaveURL('/');
  });
});

test.describe('Pipeline Editor', () => {
  test('should display React Flow canvas', async ({ page }) => {
    test.skip('Pipeline editor not implemented yet');

    await page.goto('/pipelines');

    // Vérifier que React Flow est chargé
    const reactFlowCanvas = page.locator('.react-flow');
    await expect(reactFlowCanvas).toBeVisible();
  });

  test('should add OPC-UA source node', async ({ page }) => {
    test.skip('Pipeline editor not implemented yet');

    await page.goto('/pipelines');

    // Ajouter un nœud OPC-UA
    await page.getByRole('button', { name: 'Add Source' }).click();
    await page.getByText('OPC-UA').click();

    // Vérifier que le nœud est créé
    const opcuaNode = page.locator('.ml-elec-node.opcua-source');
    await expect(opcuaNode).toBeVisible();
  });
});
