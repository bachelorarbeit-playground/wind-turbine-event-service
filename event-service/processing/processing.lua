local time = require("time")
local timeFormat = "2006-01-02 15:04:05"

function preprocess(event)
  local valueInKwh = event.value * 1000
  local timestamp = time.parse(event.date .. " " .. tostring(event.interval - 1) .. ":00:00", timeFormat, event.timezone)
  local utcTimestamp = time.format(timestamp, timeFormat, "UTC")

  local result = {
    parkId = event.parkId,
    region = event.region,
    availability = event.availability,
    value = valueInKwh,
    timestamp = utcTimestamp
  }

  return result
end

function ingestionPipeline(event)
  local preprocessedEvent = preprocess(event)
  if preprocessedEvent.availability < 30 then
    return nil
  end

  local processingTimestamp = time.format(time.unix(), timeFormat, "UTC")

  local result = {
    parkId = preprocessedEvent.parkId,
    timestamp = preprocessedEvent.timestamp,
    value = preprocessedEvent.value,
    processingTimestamp = processingTimestamp
  }

  return result
end

function anomalyDetection(event)
  local preprocessedEvent = preprocess(event)
  if preprocessedEvent.availability >= 30 then
    return nil
  end

  local regionUppercase = string.upper(preprocessedEvent.region)
  local availabilityPercentage = preprocessedEvent.availability / 100
  local processingTimestamp = time.format(time.unix(), timeFormat, "UTC")

  local result = {
    parkId = preprocessedEvent.parkId,
    timestamp = preprocessedEvent.timestamp,
    value = preprocessedEvent.value,
    region = regionUppercase,
    availability = availabilityPercentage,
    processingTimestamp = processingTimestamp
  }

  return result
end
