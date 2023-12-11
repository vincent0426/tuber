import request from './wrapper';

export class TripService {
    getHistory({ trip_status, is_driver }) {
        return request({
            method: 'get',
            url: `/trips/my?status=${trip_status}&is_driver=${is_driver}`,
            data: {}
        });
    }
    createTrip(trip) {
        return request({
            method: 'post',
            url: '/trips',
            data: trip
        });
    }
    getAllTrip() {
        return request({
            method: 'get',
            url: '/trips'
        });
    }
    getThePageTrip(PageNumber) {
        return request({
            method: 'get',
            url: '/trips?page=' + PageNumber
        });
    }
    getTrip(id) {
        return request({
            method: 'get',
            url: '/trips/' + id
        });
    }
    joinTrip(id,sourceID,destinationID) {
        return request({
            method: 'post',
            url: '/trips/' + id + '/join',
            data: {
                "source_id": sourceID,
                "destination_id": destinationID
            }
        });
    }
    getPassenger(id) {
        return request({
            method: 'get',
            url: '/trips/' + id + '/passengers'
        });
    }
    getMyTrips(is_driver){
        return request({
            method: 'get',
            url: '/trips/my?status=not_start&is_driver=' + is_driver
        });
    }
}
