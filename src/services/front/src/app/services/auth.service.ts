import { Injectable } from '@angular/core';
import { OAuthService, AuthConfig } from 'angular-oauth2-oidc';
import { BehaviorSubject, Observable } from 'rxjs';
import { ConfigService } from './config.service';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private userSubject: BehaviorSubject<Record<string, any> | null>;
  constructor(
    private oauthService: OAuthService,
    private config: ConfigService
  ) {
    oauthService.configure(this.getAuthConfig());

    oauthService.loadDiscoveryDocumentAndTryLogin().then(() => {
      if (this.oauthService.hasValidAccessToken()) {
        // User is already logged in
        console.log('LOGGEDIN');
        console.log(this.oauthService.getIdentityClaims());
        // Optional: Retrieve user profile or do any other necessary actions
      } else {
        console.log('NOT LOGGEDIN');
        
        // The user will click the login button to initiate the login process
      }
    });
    oauthService.setupAutomaticSilentRefresh();
    // Initialize the BehaviorSubject with the current user profile or null
    this.userSubject = new BehaviorSubject<Record<string, any> | null>(
      this.oauthService.getIdentityClaims() || null
    );
    // Update the BehaviorSubject whenever the user profile changes
    this.oauthService.events.subscribe((event) => {
      if (event.type === 'token_received') {
        this.userSubject.next(this.oauthService.getIdentityClaims());
      }
    });
    console.log('CurrentToken', this.userSubject.value);
  }
  public getAuthConfig(): AuthConfig {
    const appConfig = this.config.getConfig();
    return {
      // Url des Authorization-Servers
      issuer: appConfig.oidcIssuer,
      // Url der Angular-Anwendung
      // An diese URL sendet der Authorization-Server den Access Code
      redirectUri: window.location.origin + appConfig.oidcRedirectUri,
      // Name der Angular-Anwendung
      clientId: appConfig.oidcClientId,

      // Rechte des Benutzers, die die Angular-Anwendung wahrnehmen möchte
      scope: 'openid profile email offline_access',

      // Code Flow (PKCE ist standardmäßig bei Nutzung von Code Flow aktiviert)
      responseType: 'code',
      strictDiscoveryDocumentValidation: false,
    } as AuthConfig;
  }

  get userName(): string {
    const claims = this.oauthService.getIdentityClaims();
    if (!claims) return 'UNKNOWN';
    return claims['given_name'];
  }

  get idToken(): string {
    return this.oauthService.getIdToken();
  }

  get accessToken(): string {
    return this.oauthService.getAccessToken();
  }

  get identityClaims(): Record<string, any> {
    return this.oauthService.getIdentityClaims();
  }
  // Observable for the current user profile
  get currentUser$(): Observable<Record<string, any> | null> {
    return this.userSubject.asObservable();
  }
  login() {
    this.oauthService.initCodeFlow();
  }
  logout() {
    this.oauthService.logOut();
  }

  refresh() {
    this.oauthService.refreshToken();
  }
}
