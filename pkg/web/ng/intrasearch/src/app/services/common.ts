import { HttpHeaders } from '@angular/common/http';

export const collectionURL = '/vecdb/';

export const httpHeaders: HttpHeaders = new HttpHeaders({
  //Authorization: 'Bearer JWT-token'
  Accept: 'application/json',
});
